export default interface LogLine {
    lineNumber: number
    contents: string
    highlight?: {
        start: number
        end: number
    }
}
