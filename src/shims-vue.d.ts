/// <reference types="vite/client" />

declare module '*.vue' {
    import type {DefineComponent} from 'vue'
    const component: DefineComponent<object, object, any>
    export default component
}

declare module '*.jpg' {
    const value: string;
    export default value;
};