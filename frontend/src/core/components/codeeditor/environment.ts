import * as monaco from "monaco-editor"
import { loader } from "@monaco-editor/react"
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker?worker"
import CssWorker from "monaco-editor/esm/vs/language/css/css.worker?worker"
import HtmlWorker from "monaco-editor/esm/vs/language/html/html.worker?worker"
import JsonWorker from "monaco-editor/esm/vs/language/json/json.worker?worker"
import TsWorker from "monaco-editor/esm/vs/language/typescript/ts.worker?worker"

interface MonacoEnvironmentConfig {
    getWorker: (moduleId: string, label: string) => Worker
}

loader.config({ monaco })

const globalSelf = self as typeof self & { MonacoEnvironment?: MonacoEnvironmentConfig }

globalSelf.MonacoEnvironment = {
    getWorker(_moduleId, label) {
        switch (label) {
            case "json":
                return new JsonWorker()
            case "css":
            case "scss":
            case "less":
                return new CssWorker()
            case "html":
            case "handlebars":
            case "razor":
                return new HtmlWorker()
            case "typescript":
            case "javascript":
                return new TsWorker()
            default:
                return new EditorWorker()
        }
    },
}
