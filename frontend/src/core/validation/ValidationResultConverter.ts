import ApiResponse from "../apiclient/ApiResponse";
import ValidationResult from "./ValidationResult";

interface ErrorDetails {
    fieldPath: string
    message: string
}

class ValidationResultConverter {
    parse(response: ApiResponse<any>): ValidationResult | null {
        const body = response.body
        if (!Array.isArray(body?.consistencyProblems))
            return null

        const details = body.consistencyProblems as Array<ErrorDetails>
        const errors = new Map<String, Array<String>>()
        details?.forEach(errorDetails => {
            const messages = errors.get(errorDetails.fieldPath) || []
            messages.push(errorDetails.message)
            errors.set(errorDetails.fieldPath, messages)
        })

        return new ValidationResult(errors)
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new ValidationResultConverter()
