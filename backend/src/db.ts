import { Sequelize, SequelizeOptions, Model, Table, Column, DataType } from 'sequelize-typescript';

@Table
export class User extends Model<User> {
    @Column({
        type: DataType.INTEGER,
        primaryKey: true,
        allowNull: false,
        unique: true,
        autoIncrement: true,
    })
    id: number;

    @Column({ 
        type: DataType.STRING,
        allowNull: false, 
        unique: true 
    })
    name: string;

    @Column({
        type: DataType.STRING,
        allowNull: false
    })
    password: string;
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
    User
]);