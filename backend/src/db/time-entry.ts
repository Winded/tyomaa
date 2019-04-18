import { Model, Table, Column, ForeignKey, BelongsTo, DataType } from 'sequelize-typescript';
import { User } from './user';
import { nameIdentifierRegExp } from '../util/identity';

@Table
export class TimeEntry extends Model<TimeEntry> {
    @Column({
        type: DataType.INTEGER,
        primaryKey: true,
        allowNull: false,
        unique: true,
        autoIncrement: true,
    })
    id: number;

    @ForeignKey(() => User)
    @Column({
        type: DataType.INTEGER,
        allowNull: false,
    })
    userId: number;

    @Column({
        type: DataType.STRING,
        allowNull: false,
    })
    get project(): string {
        return this.getDataValue('project');
    }

    set project(value: string) {
        if (!nameIdentifierRegExp.test(value)) {
            throw new EvalError('Name must only contain alphabetic characters, numbers and dashes');
        }
        this.setDataValue('project', value);
    }

    @Column({
        type: DataType.DATE,
        allowNull: false,
    })
    start: Date;
    @Column({
        type: DataType.DATE,
        allowNull: true,
    })
    end: Date;

    @BelongsTo(() => User)
    user: User;
}