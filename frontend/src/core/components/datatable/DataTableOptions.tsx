import React from "react"
import { Modal } from "antd"

export interface DataTableOptionsProps {
    id: string
    open: boolean
    onClose: () => void
}

export default class DataTableOptions extends React.Component<DataTableOptionsProps, any> {
    constructor(props: DataTableOptionsProps) {
        super(props)
    }

    private save() {
        const { onClose } = this.props
        onClose()
    }

    private cancel() {
        const { onClose } = this.props
        onClose()
    }

    render() {
        const { open } = this.props

        return (
            <Modal title="Options" open={open} onCancel={() => this.cancel()} onOk={() => this.save()}>
                <p>TODO: Implement this</p>
            </Modal>
        )
    }
}
