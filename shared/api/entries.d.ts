import { Entry } from './entry';

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

export interface EntriesSingleResponse {
    entry: Entry;
}