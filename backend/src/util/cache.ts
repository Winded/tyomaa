import * as NodeCache from 'node-cache';

const nodeCache = new NodeCache({
    stdTTL: 60 * 15,
    deleteOnExpire: true,
});

export function get<T>(userId: number, name: string): T | undefined {
    let userCache = nodeCache.get(`user_${userId}`);
    if(!userCache) {
        return undefined;
    }

    return userCache[name];
}

export function set<T>(userId: number, name: string, value: T): boolean {
    let userCache = nodeCache.get(`user_${userId}`);
    if(!userCache) {
        userCache = {};
    }

    userCache[name] = value;
    nodeCache.set(`user_${userId}`, userCache);
    return true;
}

export function del(userId: number, name: string): boolean {
    let userCache = nodeCache.get(`user_${userId}`);
    if(!userCache) {
        return false;
    }

    if(userCache[name] === undefined) {
        return false;
    }

    userCache[name] = undefined;
    nodeCache.set(`user_${userId}`, userCache);
    return true;
}

export function clear(userId: number): void {
    nodeCache.del(`user_${userId}`);
}