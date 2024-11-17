import type {ValidateStatus} from "antd/es/form/FormItem";

export default class ValidationResult {
    private readonly errors: Map<String, Array<String>>

    constructor(errors?: Map<String, Array<String>>) {
        this.errors = errors || new Map()
    }

    getStatus(path: String): ValidateStatus | undefined {
        const errors = this.errors.get(path)
        return errors === undefined || errors.length === 0 ? undefined : "error"
    }

    getMessage(path: String): String | undefined {
        const errors = this.errors.get(path)
        return errors === undefined || errors.length === 0 ? undefined : errors.join("; ")
    }
}
