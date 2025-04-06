package services

import (
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

func GenerateLessonPlan(userContent string) (LessonPlan, error) {

	var plan LessonPlan

	ai.GeminiModel.ResponseMIMEType = "application/json"

	prompt := fmt.Sprintf(`
You are an expert educator. Create a structured learning plan based on the following text. Make sure the content is educational only. If it contains wrong information or inappropriate content respond with an error message. Refer to the response format.

Content: %s

Generate the main title for the content. preferrably a short title under 32 characters, hard limit is 64 characters.

Give me an emoji that is related to the main title of the content. Give me the UTF-8 code of the emoji

Generate 3 to 5 lessons based on the content length with the following structure: The lessons should be focused more on the theory aspect of the content.
    1. A title for the lesson (less than 64 characters)
    2. Key learning objectives (2-4 bullet points)
    3. Main content (2-3 paragraphs explaining the key concepts)
    4. 3 practice questions with answers. Questions can be multiple choice or True/False



Format the response as a JSON including a message and array of lesson objects with the following structure: The message is "success" | "error: insufficient content" | "error: <relevant error message>"
{
    "Message": message,
    "Emoji" : <UTF-8>,
    "MainTitle": <MAIN_TITLE>,
    "Lessons": [
      {
        "Title": "Lesson title",
        "Objectives": ["objective 1", "objective 2", "objective 3"],
        "Content": ["paragraph 1", "paragraph 2"],
        "Questions": [
          {
            "QuestionText": "Question text",
            "Options": ["option A", "option B", "option C", "option D"],
            "CorrectAnswer": 0
          }
        ]
      }
    ]
}
Make sure the content is educational, engaging, and follows a logical progression.
		`, userContent)

	resp, err := ai.GeminiModel.GenerateContent(ai.Ctx, genai.Text(prompt))
	if err != nil {
		log.Println(err)
		return plan, err
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Println("Received an empty or invalid response from the API.")
		return plan, err
	}

	// Iterate through parts (usually only one for JSON mode)
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			fmt.Println("\nRaw JSON Response:")
			fmt.Println(string(txt)) // Print raw JSON text

			// Unmarshal the JSON text into our Go struct
			if err := json.Unmarshal([]byte(txt), &plan); err != nil {
				log.Printf("Error unmarshalling JSON: %v\nRaw Text: %s", err, string(txt))
				return plan, err
			}

			fmt.Printf("\nAccessing data directly:\n")
			fmt.Printf("  Message: %s\n", plan.Message)
			if plan.Message != "success" {
				fmt.Printf("  Error Message: %s\n", plan.Message)
				return plan, errors.New(plan.Message)
			}

		} else {
			log.Printf("Received a part that is not text: %T\n", part)
		}
	}
	return plan, nil
}
