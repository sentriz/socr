/// <reference types="node" />

declare module 'relative-date' {
    export default (input: Date | number, reference?: Date | number) => string
}