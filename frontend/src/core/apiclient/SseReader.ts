export interface SseEvent {
    event: string
    data: string
}

const frameBoundary = "\n\n"

export async function readEventStream(response: Response, onEvent: (event: SseEvent) => void): Promise<void> {
    if (response.body === null) throw new Error("The response does not contain a body")

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ""

    try {
        while (true) {
            const { value, done } = await reader.read()
            if (done) break

            buffer += decoder.decode(value, { stream: true })
            buffer = buffer.replace(/\r\n/g, "\n")

            let boundary = buffer.indexOf(frameBoundary)
            while (boundary >= 0) {
                const block = buffer.substring(0, boundary)
                buffer = buffer.substring(boundary + frameBoundary.length)

                const event = parseEvent(block)
                if (event !== null) onEvent(event)

                boundary = buffer.indexOf(frameBoundary)
            }
        }

        buffer += decoder.decode()
        const event = parseEvent(buffer.replace(/\r\n/g, "\n"))
        if (event !== null) onEvent(event)
    } finally {
        reader.releaseLock()
    }
}

function parseEvent(block: string): SseEvent | null {
    let event = "message"
    const data: string[] = []

    for (const line of block.split("\n")) {
        if (line === "" || line.startsWith(":")) continue

        if (line.startsWith("event:")) event = line.substring("event:".length).trim()
        else if (line.startsWith("data:")) data.push(line.substring("data:".length).trimStart())
    }

    if (data.length === 0) return null

    return { event, data: data.join("\n") }
}
