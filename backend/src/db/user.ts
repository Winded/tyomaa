import { Model, Table, Column, HasMany, DataType } from 'sequelize-typescript';
import { TimeEntry } from './time-entry';
import { validateNameIdentifier } from '../util/identity';
import ApiUser from '../../../shared/api/user';
import { IApiFormattable } from '.';

@Table
export class User extends Model<User> implements IApiFormattable<ApiUser> {
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
        if (!validateNameIdentifier(value)) {
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

    toApiFormat(): ApiUser {
        return {
            id: this.id,
            name: this.name,
        };
    }
}