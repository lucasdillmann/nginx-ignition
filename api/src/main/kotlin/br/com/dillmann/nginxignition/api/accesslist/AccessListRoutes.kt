package br.com.dillmann.nginxignition.api.accesslist

import br.com.dillmann.nginxignition.api.accesslist.handler.*
import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.basePath

internal class AccessListRoutes(
    private val listHandler: ListAccessListHandler,
    private val getHandler: GetAccessListByIdHandler,
    private val putHandler: PutAccessListHandler,
    private val postHandler: PostAccessListHandler,
    private val deleteHandler: DeleteAccessListByIdHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        basePath("/api/access-lists") {
            requireAuthentication {
                get(listHandler)
                post(postHandler)

                path("/{id}") {
                    get(getHandler)
                    put(putHandler)
                    delete(deleteHandler)
                }
            }
        }
}
