package services

import (
	"akasha-api/internal/req_structs"
	"akasha-api/pkg/ai"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
)

type Question struct {
	QuestionText  string   `json:"QuestionText"`
	Options       []string `json:"Options"`
	CorrectAnswer int      `json:"CorrectAnswer"`
}

// Lesson represents a single lesson unit.
type Lesson struct {
	Title      string     `json:"Title"`
	Objectives []string   `json:"Objectives"`
	Content    []string   `json:"Content"`
	Questions  []Question `json:"Questions"`
}

// LessonPlan represents the overall structure of the response.
type LessonPlan struct {
	Message   string   `json:"Message"`
	MainTitle string   `json:"MainTitle"`
	Emoji     string   `json:"Emoji"`
	Lessons   []Lesson `json:"Lessons"`
}

func GenerateLessonPlan(req req_structs.GenerateTopicReqBody) (LessonPlan, error) {

	var plan LessonPlan

	var mode string

	if req.OnlyContent {
		mode = "strict"
	} else {
		mode = "explore"
	}

	ai.GeminiModel.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`
		You are an expert educator AI. Your task is to create a structured, theory-focused learning plan based on the provided text.

**Input Parameters:**
Content: %s
Number Of Lessons: %d
Mode: %s

**Instructions:**

1.  **Analyze Input:** First, determine if the text is suitable for creating an educational plan.
    * If the text is suitable, proceed to step 2.
    * If the text is *unsuitable* (e.g., contains factually incorrect information, harmful/inappropriate content, is fundamentally non-educational like gibberish or ads, or is too trivial/empty), **identify the primary reason** for rejection. Then, **STOP** and respond *only* with the error JSON format specified below.
2.  **Generate Title & Emoji:** If the text is suitable, generate:
    * A concise 'MainTitle' for the overall topic (preferably under 32 characters, hard limit 64 characters).
    * A relevant 'Emoji' representing the main title.
3.  **Create Lessons:** (Only if input content is suitable)
    * Generate exactly **%d** lessons based on the **Content**.
    * Adjust the depth and breadth of each lesson appropriately to fit the total number requested and the amount of source material available in the **Content**. Avoid making lessons overly thin or repetitive if a high number is requested for short content.
    * Ensure the lessons are ordered logically.
    * **Apply Generation Mode (%s):**
        * If **Mode** is "strict": Base the lesson Content paragraphs *exclusively* on the information present in the provided **Content**. Do not introduce external concepts or information not directly mentioned or clearly implied in the text.
        * If **Mode** is "explore": Base the lesson Content primarily on the provided **Content**. However, you may *enrich* the explanation slightly with closely related foundational concepts, definitions, or brief, relevant examples *only if* they directly clarify or enhance the understanding of the topics *present in the source text*. Ensure the core focus remains tightly bound to the input text's subject matter and avoid introducing unrelated topics.
    * Each lesson must include:
        * 'Title': A short title (less than 64 characters).
        * 'Objectives': 2-4 bullet points listing key learning goals for the lesson.
        * 'Content': 2-3 paragraphs explaining the key theoretical concepts from the text related to the objectives. Use clear, concise language. Where appropriate and supported by the source text, briefly mention real-world relevance or examples to enhance engagement.
        * 'Questions': 3 practice questions (Multiple Choice or True/False) with answers.
            * Ensure each question directly assesses understanding of the key concepts or objectives presented *in that specific lesson's content*.
            * For True/False questions, use 'Options: ["True", "False"]' and 'CorrectAnswer: 0' for True, '1' for False.
            * For Multiple Choice, provide 3-4 plausible options including the correct one.

4.  **Format Output:** Format the entire response as a single JSON object adhering *strictly* to the structure below.

**JSON Output Format (Success):**
{
  "Message": "success",
  "Emoji": "<GENERATED_EMOJI>",
  "MainTitle": "<GENERATED_MAIN_TITLE>",
  "Lessons": [
    {
      "Title": "Lesson 1 title",
      "Objectives": ["objective 1.1", "objective 1.2"],
      "Content": ["paragraph 1 explaining concepts...", "paragraph 2 expanding on concepts..."],
      "Questions": [
        {
          "QuestionText": "Question 1 text?",
          "Options": ["Option A", "Option B", "Option C"],
          "CorrectAnswer": 0 // Index of the correct option
        },
        {
          "QuestionText": "True or False: Statement?",
          "Options": ["True", "False"],
          "CorrectAnswer": 1 // Index 1 = False
        },
        // ... more questions
      ]
    },
    // ... more lessons (total 3-5)
  ]
}

**JSON Output Format (Error):**
{
  "Message": "<BEGIN_WITH_req_error:_THEN_AI_GENERATES_ERROR_MESSAGE>",
  "Emoji": null,
  "MainTitle": null,
  "Lessons": []
}

**Instructions for Generating the Error Message (if input is unsuitable):**
* The 'Message' field in the error JSON must start *exactly* with 'req_error: '.
* Following 'req_error: ', generate a *unique and creative* error message.
* This message should be:
    * Fun and informal in tone.
    * Include a relevant emoji (related to the error type if possible).
    * Briefly and playfully *hint* at the reason the input was rejected (e.g., "This looks like gibberish!", "Can't make lessons from that!", "Needs more educational spice!").
    * Strictly less than 20 words (including the emoji).

**(Self-Correction Note for AI):** Remember to first analyze *why* the content is unsuitable before crafting the dynamic error message for the 'Message' field in the error JSON. Adhere strictly to all constraints for the error message (prefix, tone, length, emoji, hinting at reason).
		`, req.Content, req.NumOfLessons, mode, req.NumOfLessons, mode)

	resp, err := ai.GeminiModel.GenerateContent(ai.Ctx, genai.Text(prompt))
	if err != nil {
		log.Println(err)
		return plan, err
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("Received an empty or invalid response from the API.")
		return plan, err
	}

	// Log the raw response
	// log.Println(resp.Candidates[0].Content.Parts[0])

	// Iterate through parts (usually only one for JSON mode)
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {

			// Unmarshal the JSON text into our Go struct
			if err := json.Unmarshal([]byte(txt), &plan); err != nil {
				log.Printf("Error unmarshalling JSON: %v\nRaw Text: %s", err, string(txt))
				return plan, err
			}

			if plan.Message != "success" {
				log.Println(req.Content)
				log.Printf("  Error Message: %s\n", plan.Message)
				return plan, errors.New(plan.Message)
			}

		} else {
			log.Printf("Received a part that is not text: %T\n", part)
		}
	}
	return plan, nil
}
