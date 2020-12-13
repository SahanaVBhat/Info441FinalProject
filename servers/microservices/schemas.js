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
    description: {type: String},
    likedUsers: {type:[{studentID:Number}]},
    dislikedUsers: {type:[{studentID:Number}]},
    createdAt: { type: Date, required: true},
    editedAt: {type: Date}
});


const course = new Schema({
	id: {type:Number, required:true, unique:true},
    code: {type: String, require: true},
    title: {type: String}, 
    description: {type: String}, 
    credits: {type: Number}
});

module.exports = {
    evaluation, course
};