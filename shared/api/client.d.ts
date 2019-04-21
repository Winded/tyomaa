import { TokenPostResponse, TokenPostRequest, TokenGetResponse } from "./auth";
import { ClockGetResponse, ClockStartPostRequest, ClockStartPostResponse, ClockStopPostResponse } from "./clock";
import { EntriesGetResponse, EntriesGetRequest, EntriesPostRequest, EntriesPostResponse, EntriesSingleGetResponse, EntriesSinglePostRequest, EntriesSinglePostResponse } from "./entries";
import { ProjectsGetResponse } from './projects';

export interface IApiClient {
    token: string;

    authTokenGet(): Promise<TokenGetResponse>;
    authTokenPost(body: TokenPostRequest): Promise<TokenPostResponse>;

    clockGet(): Promise<ClockGetResponse>;
    clockStartPost(body: ClockStartPostRequest): Promise<ClockStartPostResponse>;
    clockStopPost(): Promise<ClockStopPostResponse>;

    entriesGet(query: EntriesGetRequest): Promise<EntriesGetResponse>;
    entriesPost(body: EntriesPostRequest): Promise<EntriesPostResponse>;
    entriesSingleGet(entryId: number): Promise<EntriesSingleGetResponse>;
    entriesSinglePost(entryId: number, body: EntriesSinglePostRequest): Promise<EntriesSinglePostResponse>;

    projectsGet(): Promise<ProjectsGetResponse>;
}