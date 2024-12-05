package br.com.dillmann.nginxignition.database.common

import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.like

fun Query.withSearchTerms(searchTerms: String?, fields: List<Column<out Any>>): Query {
    if (searchTerms.isNullOrBlank() || fields.isEmpty()) return this

    val pattern = searchTerms.replace(" ", "%").uppercase().let { "%$it%" }
    val predicates = fields.map { it.castTo(VarCharColumnType()).upperCase() like pattern }
    return where { OrOp(predicates) }
}
