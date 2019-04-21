export interface Project {
    name: string;
    totalTime: number;
}

export interface ProjectsGetResponse {
    projects: Project[];
}