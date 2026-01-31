import React from "react"
import { Modal } from "antd"
import DataTablePersistentStateConfig from "./model/DataTablePersistentStateConfig"
import DataTableService from "./DataTableService"

export interface DataTableOptionsProps {
    id: string
    open: boolean
    onClose: () => void
}

interface DataTableOptionsState {
    config: DataTablePersistentStateConfig
}

export default class DataTableOptions extends React.Component<DataTableOptionsProps, DataTableOptionsState> {
    private readonly service: DataTableService

    constructor(props: DataTableOptionsProps) {
        super(props)
        this.service = new DataTableService()

        this.state = {
            config: this.service.currentConfig(),
        }
    }

    private save() {
        const { config } = this.state
        this.service.updateConfig(config)

        const { onClose } = this.props
        onClose()
    }

    private cancel() {
        const { onClose } = this.props
        onClose()
    }

    componentDidUpdate(prevProps: Readonly<DataTableOptionsProps>) {
        const { open } = this.props

        if (prevProps.open !== open) {
            this.setState({
                config: this.service.currentConfig(),
            })
        }
    }

    render() {
        const { open } = this.props
        const { config } = this.state

        return (
            <Modal title="Options" open={open} onCancel={() => this.cancel()} onOk={() => this.save()}>
                <p>TODO: Implement this</p>
            </Modal>
        )
    }
}
