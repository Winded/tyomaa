import { Sequelize, SequelizeOptions } from 'sequelize-typescript';
import { User } from './user';
import { TimeEntry } from './time-entry';

export interface IApiFormattable<T> {
    toApiFormat(): T;
}

export const sequelize = new Sequelize({
    dialect: process.env.DB_IN_MEMORY ? 'sqlite' : 'postgres',
    storage: process.env.DB_IN_MEMORY ? ':memory:' : undefined,
    host: process.env.DB_HOST || 'localhost',
    port: process.env.DB_PORT || 5432,
    database: process.env.DB_DATABASE || 'tyomaa',
    username: process.env.DB_USER || 'tyomaa',
    password: process.env.DB_PASSWORD || 'tyomaa',
} as SequelizeOptions);

sequelize.addModels([
    User,
    TimeEntry,
]);