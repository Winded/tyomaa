import { Entry } from './entry';

export interface EntriesGetRequest {
    // TODO
}
export interface EntriesGetResponse {
    entries: Entry[];
}

export interface EntriesPostRequest {
    project: string;
    start: string;
    end: string;
}

export interface EntriesPostResponse {
    entry: Entry;
}

export interface EntriesSingleGetResponse {
    entry: Entry;
}

export interface EntriesSinglePostRequest extends EntriesPostRequest {}
export interface EntriesSinglePostResponse extends EntriesPostResponse {}