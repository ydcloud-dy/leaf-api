import { PropType } from 'vue';
export type IconName = 'bold' | 'underline' | 'italic' | 'strike-through' | 'title' | 'sub' | 'sup' | 'quote' | 'unordered-list' | 'ordered-list' | 'task' | 'code-row' | 'code' | 'link' | 'image' | 'table' | 'revoke' | 'next' | 'save' | 'prettier' | 'minimize' | 'maximize' | 'fullscreen-exit' | 'fullscreen' | 'preview-only' | 'preview' | 'preview-html' | 'catalog' | 'github' | 'mermaid' | 'formula' | 'close' | 'delete' | 'upload';
declare const _default: import("vue").DefineComponent<import("vue").ExtractPropTypes<{
    name: {
        type: PropType<IconName>;
        default: string;
    };
}>, () => import("vue").VNode<import("vue").RendererNode, import("vue").RendererElement, {
    [key: string]: any;
}>, {}, {}, {}, import("vue").ComponentOptionsMixin, import("vue").ComponentOptionsMixin, {}, string, import("vue").PublicProps, Readonly<import("vue").ExtractPropTypes<{
    name: {
        type: PropType<IconName>;
        default: string;
    };
}>> & Readonly<{}>, {
    name: IconName;
}, {}, {}, {}, string, import("vue").ComponentProvideOptions, true, {}, any>;
export default _default;
