export default function resolveLanguageTag(
    availableLanguages: string[],
    currentLanguage: string | null,
    defaultLanguage: string,
): string {
    const targetLanguages = (currentLanguage ?? defaultLanguage).split(",")

    for (const targetLanguage of targetLanguages) {
        if (availableLanguages.includes(targetLanguage)) {
            return targetLanguage
        }

        if (targetLanguage.includes("-")) {
            const baseLanguage = targetLanguage.split("-")[0]
            if (availableLanguages.includes(baseLanguage)) {
                return baseLanguage
            }
        }
    }

    const defaultBaseLanguage = defaultLanguage.split("-")[0]
    const fallback = availableLanguages.find(
        language => language === defaultLanguage || language === defaultBaseLanguage,
    )

    return fallback ?? defaultLanguage
}
