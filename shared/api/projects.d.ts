export interface Project {
    name: string;
    totalTime: string;
}

export interface ProjectsGetResponse {
    projects: Project[];
}