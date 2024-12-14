export function buildLoginUrl() {
    // eslint-disable-next-line no-restricted-globals
    const { pathname, search } = location
    const returnTo = !search ? pathname : `${pathname}?${search}`
    return `/login?returnTo=${encodeURI(returnTo)}`
}
