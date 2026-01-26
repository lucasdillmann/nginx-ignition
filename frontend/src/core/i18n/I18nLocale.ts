import I18nContext from "./I18nContext"
import MessageKey from "./model/MessageKey.generated"
import { i18n } from "./I18n"

export function buildI18nLocale() {
    const { currentLanguage, defaultLanguage } = I18nContext.get()

    return {
        locale: currentLanguage ?? defaultLanguage,
        Form: {
            optional: `(${i18n(MessageKey.CommonOptional)})`,
            defaultValidateMessages: {},
        },
        Empty: {
            description: i18n(MessageKey.CommonNoData),
        },
        Pagination: {
            items_per_page: i18n(MessageKey.FrontendLocalePaginationItemsPerPage),
            jump_to: i18n(MessageKey.FrontendLocalePaginationJumpTo),
            jump_to_confirm: i18n(MessageKey.FrontendLocalePaginationJumpToConfirm),
            page: i18n(MessageKey.FrontendLocalePaginationPage),
            prev_page: i18n(MessageKey.FrontendLocalePaginationPrevPage),
            next_page: i18n(MessageKey.FrontendLocalePaginationNextPage),
            prev_5: i18n(MessageKey.FrontendLocalePaginationPrev5),
            next_5: i18n(MessageKey.FrontendLocalePaginationNext5),
            prev_3: i18n(MessageKey.FrontendLocalePaginationPrev3),
            next_3: i18n(MessageKey.FrontendLocalePaginationNext3),
            page_size: i18n(MessageKey.FrontendLocalePaginationPageSize),
        },
    }
}
