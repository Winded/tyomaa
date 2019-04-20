import { Entry } from './entry';

export interface EntriesGetResponse {
    entries: Entry[];
}

export interface EntriesPostRequest {
    project: string;
    start: Date;
    end: Date;
}

export interface EntriesPostResponse {
    entry: Entry;
}

export interface EntriesSingleResponse {
    entry: Entry;
}