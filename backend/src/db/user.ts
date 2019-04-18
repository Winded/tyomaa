import { Model, Table, Column, HasMany, DataType } from 'sequelize-typescript';
import { TimeEntry } from './time-entry';
import { nameIdentifierRegExp } from '../util/identity';

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
    get name(): string {
        return this.getDataValue('name');
    }

    set name(value: string) {
        if (!nameIdentifierRegExp.test(value)) {
            throw new EvalError('Name must only contain alphabetic characters, numbers and dashes');
        }
        this.setDataValue('name', value);
    }

    @Column({
        type: DataType.STRING,
        allowNull: false
    })
    password: string;

    @HasMany(() => TimeEntry)
    entries: TimeEntry[];
}