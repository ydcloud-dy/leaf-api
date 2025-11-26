import CodeMirrorUt from '../layouts/Content/codemirror';
export type ToolDirective = 'bold' | 'underline' | 'italic' | 'strikeThrough' | 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6' | 'quote' | 'unorderedList' | 'orderedList' | 'task' | 'codeRow' | 'code' | 'link' | 'image' | 'table' | 'sub' | 'sup' | 'prettier' | 'flow' | 'sequence' | 'gantt' | 'class' | 'state' | 'pie' | 'relationship' | 'journey' | 'katexInline' | 'katexBlock' | 'universal';
/**
 *
 * @param direct 操作指令
 * @param codeMirrorUt 编辑区辅助实例
 * @param params 自定义参数
 *
 * @returns string
 */
export declare const directive2flag: (direct: ToolDirective, codeMirrorUt: CodeMirrorUt, params?: any) => Promise<{
    text: string;
    options: {
        deviationStart: number;
        replaceStart: number;
        replaceEnd: number;
    };
} | {
    text: any;
    options: {
        select: boolean;
        replaceAll: boolean;
    };
} | {
    text: string;
    options: {
        deviationStart: number;
        deviationEnd: number;
    };
} | {
    text: string;
    options: {
        select: boolean;
        deviationStart: number;
        deviationEnd: number;
    };
} | {
    text: string;
    options: {
        select?: undefined;
        deviationStart?: undefined;
        deviationEnd?: undefined;
    };
}>;
