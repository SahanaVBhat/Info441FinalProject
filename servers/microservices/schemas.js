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


const message = new Schema({
	id: {type:Number, required:true, unique:true},
    channelID: {type:Number, required:true},
    body: {type:String, required:true},
    createdAt: { type: Date, required: true},
    creator: {type:{id:Number, email:String}},
    editedAt: Date
});

module.exports = {
    channel,
    message
};