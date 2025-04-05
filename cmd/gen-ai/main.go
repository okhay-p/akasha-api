package main

import (
	"akasha-api/pkg/config"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Question struct {
	QuestionText  string   `json:"Question"` // Field name matches "Question" in schema
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
	Message              string   `json:"Message"`                        // Should be "success" or "error"
	RelevantErrorMessage *string  `json:"RelevantErrorMessage,omitempty"` // Pointer for nullable string
	MainTitle            string   `json:"MainTitle"`
	Lessons              []Lesson `json:"Lessons"`
}

func main() {

	ctx := context.Background()
	cfg := config.LoadConfig()
	apiKey := cfg.GeminiApiKey
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	model := client.GenerativeModel("gemini-2.0-flash")

	model.ResponseMIMEType = "application/json"

	model.ResponseSchema = &genai.Schema{
		// Required: []string{"Message, MainTitle, Lessons"},
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"Message": {
				Type: genai.TypeString,
				Enum: []string{"success", "error"},
			},
			"RelevantErrorMessage": {
				Type:        genai.TypeString,
				Nullable:    true,
				Description: "To indicate the reason for error",
			},
			"MainTitle": {
				Type: genai.TypeString,
			},
			"Lessons": {
				Nullable: false,
				Type:     genai.TypeArray,
				Items: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"Title": {Type: genai.TypeString},
						"Objectives": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type:        genai.TypeString,
								Description: "A learning objective for the lesson",
							},
							Description: "A list of learning objectives for the lesson",
						},
						"Content": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type:        genai.TypeString,
								Description: "A paragraph of content for the lesson",
							},
							Description: "A list of paragraphs of content for the lesson",
						},
						"Questions": {
							Type: genai.TypeArray,
							Items: &genai.Schema{
								Type: genai.TypeObject,
								Properties: map[string]*genai.Schema{
									"Question": {
										Type:        genai.TypeString,
										Description: "The text of the question.",
									},
									"Options": {
										Type: genai.TypeArray,
										Items: &genai.Schema{
											Type: genai.TypeString,
										},
									},
									"CorrectAnswer": {
										Type: genai.TypeInteger,
									},
								}, // Question Items Properties
							}, // Question Items
							Description: "A list of questions for the lesson.",
						}, // Questions
					}, // Items Properties
				}, // Lesson Items
			}, // Lessons
		}, // Base Properties
	}

	userContent :=
		`
		The Euclidean algorithm is a method for finding the greatest common divisor (GCD) of two integers. The GCD of two numbers is the largest positive integer that divides both numbers without leaving a remainder. The algorithm is based on the principle that the GCD of two numbers also divides their difference.

`

	prompt := fmt.Sprintf(
		`You are an expert educator. Create a structured learning plan based on the following text. Make sure the content is educational only. If it contains wrong information or inappropriate content respond with an error message. Refer to the response format.

Content: %s

Generate a lesson plan using the provided JSON schema. You MUST include at least one lesson in the 'Lessons' array.
    1. A title for the lesson (less than 64 characters)
    2. Key learning objectives (2-4 bullet points)
    3. Main content (2-3 paragraphs explaining the key concepts)
    4. 3 practice questions with answers

Generate the main title for the content. title should be under 64 characters, preferrably a short one.
Make sure the content is educational, engaging, and follows a logical progression.
`,
		userContent)

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		log.Fatal(err)
	}

	if resp == nil || len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		log.Fatal("Received an empty or invalid response from the API.")
		return
	}

	// Iterate through parts (usually only one for JSON mode)
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			fmt.Println("\nRaw JSON Response:")
			fmt.Println(string(txt)) // Print raw JSON text

			// Declare a variable of our target struct type
			var plan LessonPlan

			// Unmarshal the JSON text into our Go struct
			if err := json.Unmarshal([]byte(txt), &plan); err != nil {
				log.Fatalf("Error unmarshalling JSON: %v\nRaw Text: %s", err, string(txt))
			}

			// --- 7. Print the Parsed Data (Pretty Printed) ---
			fmt.Println("\nParsed Go Struct:")
			prettyJSON, err := json.MarshalIndent(plan, "", "  ") // Indent with two spaces
			if err != nil {
				log.Fatalf("Error formatting output JSON: %v", err)
			}
			fmt.Println(string(prettyJSON))

			// You can now access the data directly via the struct fields:
			fmt.Printf("\nAccessing data directly:\n")
			fmt.Printf("  Message: %s\n", plan.Message)
			if plan.Message == "error" && plan.RelevantErrorMessage != nil {
				fmt.Printf("  Error Message: %s\n", *plan.RelevantErrorMessage)
			}
			fmt.Printf("  Main Title: %s\n", plan.MainTitle)
			if len(plan.Lessons) > 0 {
				fmt.Printf("  First Lesson Title: %s\n", plan.Lessons[0].Title)
				if len(plan.Lessons[0].Questions) > 0 {
					fmt.Printf("  First Question of First Lesson: %s\n", plan.Lessons[0].Questions[0].QuestionText)
				}
			}

		} else {
			log.Printf("Received a part that is not text: %T\n", part)
		}
	}
}
