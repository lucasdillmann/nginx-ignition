import React from "react"
import { Flex } from "antd"
import "./LogViewer.css"
import LogLine from "../model/LogLine"

export interface LogViewerProps {
    lines: LogLine[]
}

export default class LogViewer extends React.Component<LogViewerProps> {
    private readonly containerRef: React.RefObject<HTMLDivElement | null>

    constructor(props: LogViewerProps) {
        super(props)
        this.containerRef = React.createRef()
    }

    componentDidUpdate() {
        if (this.containerRef.current) this.containerRef.current.scrollTop = this.containerRef.current.scrollHeight
    }

    private renderLineContent(line: LogLine) {
        const { contents, highlight } = line

        if (!highlight) {
            return <span className="log-viewer-line-text">{contents}</span>
        }

        const { start, end } = highlight
        const before = contents.substring(0, start)
        const highlighted = contents.substring(start, end + 1)
        const after = contents.substring(end + 1)

        return (
            <span className="log-viewer-line-text">
                {before}
                <mark className="log-viewer-highlight">{highlighted}</mark>
                {after}
            </span>
        )
    }

    render() {
        const { lines } = this.props
        const sortedLines = [...lines].sort((left, right) => left.lineNumber - right.lineNumber)
        const maxLineNumber = sortedLines.length > 0 ? sortedLines[sortedLines.length - 1].lineNumber : 0
        const lineNumberWidth = Math.max(String(maxLineNumber).length, 3)

        return (
            <Flex ref={this.containerRef} className="log-viewer-container" vertical>
                <Flex className="log-viewer-content" vertical>
                    {sortedLines.map(line => (
                        <Flex key={line.lineNumber} className="log-viewer-line">
                            <span className="log-viewer-line-number" style={{ minWidth: `${lineNumberWidth}ch` }}>
                                {line.lineNumber}
                            </span>
                            {this.renderLineContent(line)}
                        </Flex>
                    ))}
                </Flex>
            </Flex>
        )
    }
}
