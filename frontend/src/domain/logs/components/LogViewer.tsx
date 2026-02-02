import React from "react"
import { Flex } from "antd"
import "./LogViewer.css"
import LogLine from "../model/LogLine"

export interface LogViewerProps {
    lines: LogLine[]
}

interface LogViewerState {
    sortedLines: LogLine[]
}

export default class LogViewer extends React.Component<LogViewerProps, LogViewerState> {
    private readonly containerRef: React.RefObject<HTMLDivElement | null>

    constructor(props: LogViewerProps) {
        super(props)
        this.containerRef = React.createRef()
        this.state = {
            sortedLines: this.sortLines(props.lines),
        }
    }

    private sortLines(lines: LogLine[]): LogLine[] {
        return [...lines].sort((left, right) => left.lineNumber - right.lineNumber)
    }

    componentDidUpdate(prevProps: Readonly<LogViewerProps>): void {
        if (this.containerRef.current) this.containerRef.current.scrollTop = this.containerRef.current.scrollHeight

        if (prevProps.lines !== this.props.lines) {
            const { lines } = this.props
            const sortedLines = this.sortLines(lines)
            this.setState({ sortedLines })
        }
    }

    private renderLineContent(line: LogLine) {
        const { contents, highlight } = line

        if (!highlight) {
            return <span className="log-viewer-line-text">{contents}</span>
        }

        const { start, end } = highlight
        const before = contents.substring(0, start)
        const highlighted = contents.substring(start, end)
        const after = contents.substring(end)

        return (
            <span className="log-viewer-line-text">
                {before}
                <mark className="log-viewer-highlight">{highlighted}</mark>
                {after}
            </span>
        )
    }

    private renderGapIndicator(lineNumberWidth: number) {
        return (
            <Flex className="log-viewer-line log-viewer-gap">
                <span className="log-viewer-line-number" style={{ width: `${lineNumberWidth}ch` }}>
                    ...
                </span>
                <span className="log-viewer-line-text log-viewer-gap-text"></span>
            </Flex>
        )
    }

    private renderLines(sortedLines: LogLine[], lineNumberWidth: number) {
        const elements: React.ReactNode[] = []

        sortedLines.forEach((line, index) => {
            if (index > 0) {
                const previousLineNumber = sortedLines[index - 1].lineNumber
                if (line.lineNumber > previousLineNumber + 1) {
                    elements.push(
                        <React.Fragment key={`gap-${previousLineNumber}-${line.lineNumber}`}>
                            {this.renderGapIndicator(lineNumberWidth)}
                        </React.Fragment>,
                    )
                }
            }

            elements.push(
                <Flex key={line.lineNumber} className="log-viewer-line">
                    <span className="log-viewer-line-number" style={{ width: `${lineNumberWidth}ch` }}>
                        {line.lineNumber + 1}
                    </span>
                    {this.renderLineContent(line)}
                </Flex>,
            )
        })

        return elements
    }

    render() {
        const { sortedLines } = this.state
        const maxLineNumber = sortedLines.length > 0 ? sortedLines[sortedLines.length - 1].lineNumber : 0
        const lineNumberWidth = Math.max(String(maxLineNumber).length, 6)

        return (
            <Flex ref={this.containerRef} className="log-viewer-container" vertical>
                <Flex className="log-viewer-content" vertical>
                    {this.renderLines(sortedLines, lineNumberWidth)}
                </Flex>
            </Flex>
        )
    }
}
