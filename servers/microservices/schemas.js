const mongoose = require('mongoose');
const Schema = require('mongoose').Schema;

const evaluation = new Schema({
    id: { type: Number, required: true, unique: true},
    studentID: {type:Number},
    courseID: {type:Number},
    instructors: {type:[{name:String}]},
    year: {type: Date},
    quarter: {type: String},
    creditType: {type: String},
    credit: {type: Number},
    workload: {type: Number},
    gradingTechniques: {type: Number},
    description: {type: string},
    likedUsers: {type:[{studentID:Number}]},
    dislikedUsers: {type:[{studentID:Number}]},
    createdAt: { type: Date, required: true},
    editedAt: Date
});


const course = new Schema({
	id: {type:Number, required:true, unique:true},
    code: {type: String, require: true},
    title: String, 
    description: String, 
    credits: Number
});

module.exports = {
    evaluation, course
};