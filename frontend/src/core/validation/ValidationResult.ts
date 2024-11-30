import type { ValidateStatus } from "antd/es/form/FormItem"

export default class ValidationResult {
    private readonly errors: Map<string, string[]>

    constructor(errors?: Map<string, string[]>) {
        this.errors = errors || new Map()
    }

    getStatus(path: string): ValidateStatus | undefined {
        const errors = this.errors.get(path)
        return errors === undefined || errors.length === 0 ? undefined : "error"
    }

    getMessage(path: string): string | undefined {
        const errors = this.errors.get(path)
        return errors === undefined || errors.length === 0 ? undefined : errors.join("; ")
    }
}
