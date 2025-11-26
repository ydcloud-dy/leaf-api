/**
 * 清空组件带来的副作用，例如
 * 1. 使用CDN嵌入的链接，为了保证多个组件能够正常使用，组件在卸载时不会主动移除
 */
export declare const clearSideEffects: () => void;
