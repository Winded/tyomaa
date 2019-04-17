import { Model, Table, Column, ForeignKey, BelongsTo, DataType } from 'sequelize-typescript';
import { User } from './user';

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
    project: string;

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