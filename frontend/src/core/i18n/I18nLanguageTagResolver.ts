export default function resolveLanguageTag(
    availableLanguages: string[],
    currentLanguage: string | null,
    defaultLanguage: string,
): string {
    const targetLanguages = (currentLanguage ?? defaultLanguage).split(",")

    for (const targetLanguage of targetLanguages) {
        for (const language of availableLanguages) {
            if (language === targetLanguage) {
                return language
            }
        }

        if (targetLanguage.includes("-")) {
            const baseLanguage = targetLanguage.split("-")[0]
            for (const language of availableLanguages) {
                if (language === baseLanguage) {
                    return language
                }
            }
        }
    }

    const defaultBaseLanguage = defaultLanguage.split("-")[0]
    for (const language of availableLanguages) {
        if (language === defaultLanguage || language === defaultBaseLanguage) {
            return language
        }
    }

    return defaultLanguage
}
