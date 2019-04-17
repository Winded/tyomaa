import { Model, Table, Column, HasMany, DataType } from 'sequelize-typescript';
import { TimeEntry } from './time-entry';

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

    @HasMany(() => TimeEntry)
    entries: TimeEntry[];
}