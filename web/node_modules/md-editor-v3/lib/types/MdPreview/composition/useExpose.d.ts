import type { SetupContext } from 'vue';
import { EditorEmits, MdPreviewProps } from '../../MdEditor/type';
export declare const useExpose: (props: MdPreviewProps, ctx: SetupContext<EditorEmits>, options: {
    editorId: string;
}) => void;
