import { ContentPreviewProps } from '../ContentPreview';
/**
 * 注册katex扩展到页面
 *
 */
declare const useMermaid: (props: ContentPreviewProps) => {
    reRenderRef: import("vue").ShallowRef<number, number>;
    replaceMermaid: () => Promise<void>;
};
export default useMermaid;
