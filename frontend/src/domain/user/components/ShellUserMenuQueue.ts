export enum QueueAction {
    OPEN_TOTP_CONFIG,
}

class ShellUserMenuQueue {
    private readonly queue: QueueAction[] = []
    private listener: ((action: QueueAction) => void) | null = null

    attach(listener: (action: QueueAction) => void) {
        this.listener = listener
        this.processQueue()
    }

    detach() {
        this.listener = null
    }

    openTotpConfig() {
        this.queue.push(QueueAction.OPEN_TOTP_CONFIG)
        this.processQueue()
    }

    private processQueue() {
        if (!this.listener) {
            return
        }

        while (this.queue.length > 0) {
            const action = this.queue.shift()
            if (action !== undefined) {
                this.listener(action)
            }
        }
    }
}

export default new ShellUserMenuQueue()
