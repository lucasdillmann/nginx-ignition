import React from "react"
import { Statistic } from "antd"
import { Link } from "react-router-dom"
import { I18n, I18nMessage } from "../../../core/i18n/I18n"

export interface CountCardProps {
    title: I18nMessage
    count: number
    linkTo: string
}

export default class CountCard extends React.PureComponent<CountCardProps> {
    render() {
        const { title, count, linkTo } = this.props

        return (
            <Link to={linkTo} className="traffic-stats-stat-card-link">
                <div className="traffic-stats-stat-card">
                    <Statistic title={<I18n id={title} />} value={count} />
                </div>
            </Link>
        )
    }
}
