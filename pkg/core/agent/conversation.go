package agent

import (
	"fmt"
	"strings"

	"google.golang.org/genai"
)

func (a *Agent) InvokeAgent(prompt string) {
	currentPrompt := prompt
	var pendingFunctionResponse *genai.FunctionResponse

	for turn := 1; turn <= a.maxTurns; turn++ {
		response, err := a.sendMessage(currentPrompt, pendingFunctionResponse)
		if err != nil {
			a.ui.PrintError(fmt.Errorf("error in turn %d: %w", turn, err))
			continue
		}

		currentPrompt = ""
		pendingFunctionResponse = nil

		nextPrompt, functionResponse, shouldContinue := a.processResponse(response)
		if !shouldContinue {
			return
		}

		currentPrompt = nextPrompt
		pendingFunctionResponse = functionResponse
	}

	a.ui.Print("Maximum number of turns reached. Ending conversation.")
}

func (a *Agent) sendMessage(prompt string, functionResponse *genai.FunctionResponse) (*genai.GenerateContentResponse, error) {
	var parts []genai.Part

	if functionResponse != nil {
		parts = append(parts, genai.Part{FunctionResponse: functionResponse})
	}

	if prompt != "" {
		parts = append(parts, genai.Part{Text: prompt})
	}

	response, err := a.chatSession.SendMessage(a.context, parts...)
	if err != nil {
		return nil, fmt.Errorf("error sending message: %w", err)
	}

	return response, nil
}

func (a *Agent) processResponse(response *genai.GenerateContentResponse) (string, *genai.FunctionResponse, bool) {
	if len(response.Candidates) == 0 {
		a.ui.Print("No response from the model. Trying again...")
		return "", nil, true
	}

	candidate := response.Candidates[0]
	if len(candidate.Content.Parts) == 0 {
		a.ui.Print("No content parts in response. Trying again...")
		return "", nil, true
	}

	part := candidate.Content.Parts[0]

	if textPart := part.Text; textPart != "" {
		return a.handleTextResponse(textPart)
	}

	if funcCall := part.FunctionCall; funcCall != nil {
		return a.handleFunctionCall(funcCall)
	}

	a.ui.Print("Unexpected response format. Trying again...")
	return "", nil, true
}

func (a *Agent) handleTextResponse(text string) (string, *genai.FunctionResponse, bool) {
	a.ui.PrintAgent(text)

	if a.isUserInputRequested(text) {
		userInput, err := a.getUserInput()
		if err != nil {
			a.ui.PrintError(fmt.Errorf("error during user interaction: %w", err))
			return "", nil, true
		}
		return GetUserInputPrompt(userInput), nil, true
	}

	textLower := strings.ToLower(text)
	if (strings.Contains(textLower, "i will change") ||
		strings.Contains(textLower, "i will modify") ||
		strings.Contains(textLower, "i will update") ||
		strings.Contains(textLower, "has been changed") ||
		strings.Contains(textLower, "has been modified") ||
		strings.Contains(textLower, "has been updated")) &&
		!strings.Contains(text, "**Conclusion:**") {

		promptForAction := "You mentioned making changes but didn't use the writeFile function. You MUST use writeFile to actually implement the changes. Please call writeFile now with the modified content."
		return promptForAction, nil, true
	}

	if strings.Contains(text, "**Conclusion:**") {
		return "", nil, false
	}

	return "", nil, true
}

func (a *Agent) handleFunctionCall(funcCall *genai.FunctionCall) (string, *genai.FunctionResponse, bool) {
	switch funcCall.Name {
	case "readFile":
		a.ui.PrintReadTool(funcCall.Args)
	case "writeFile":
		a.ui.PrintWriteTool(funcCall.Args)
	case "shellExecute":
		a.ui.PrintShellTool(funcCall.Args)
	default:
		a.ui.PrintTool(funcCall.Name, funcCall.Args)
	}

	confirmed, err := a.confirmExecution()
	if err != nil {
		a.ui.PrintError(fmt.Errorf("error during confirmation: %w", err))
		return "Could not run the tool. Try again.", nil, true
	}

	if !confirmed {
		a.ui.Print("Function execution cancelled by user.")
		return GetPermissionDeniedPrompt(funcCall.Name, fmt.Sprintf("%v", funcCall.Args)), nil, false
	}

	output, err := a.executeFunctionCall(funcCall)
	if err != nil {
		a.ui.PrintError(fmt.Errorf("error executing function call: %w", err))
		errorPrompt := fmt.Sprintf("The function call failed with error: %v. Please provide an alternative solution.", err)
		return errorPrompt, nil, true
	}

	a.ui.PrintSuccess("Function successfully executed")

	functionResponse := &genai.FunctionResponse{
		Name:     funcCall.Name,
		Response: map[string]interface{}{"result": output},
	}

	return "", functionResponse, true
}

func (a *Agent) getUserInput() (string, error) {
	a.ui.Print("User Interaction Required:")
	return a.ui.Input("Please provide any necessary suggestion or leave blank: ")
}

func (a *Agent) confirmExecution() (bool, error) {
	return a.ui.Confirm("Do you want to proceed?")
}

func (a *Agent) isUserInputRequested(text string) bool {
	return strings.Contains(strings.ToLower(text), UserInputSignature)
}
