export const nameIdentifierRegExp = /^[a-zA-Z0-9\-\/]+$/;

export function validateNameIdentifier(value: string): boolean {
    return nameIdentifierRegExp.test(value);
}