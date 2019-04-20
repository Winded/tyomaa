import User from './user';

export interface TokenGetResponse {
    token?: string;
    user?: User;
}

export interface TokenPostRequest {
    username: string;
    password: string;
}

export interface TokenPostResponse {
    token: string;
    message: string;
}