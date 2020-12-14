import React, { Component } from 'react';
import api from '../../../../Constants/APIEndpoints/APIEndpoints';
import Errors from '../../../Errors/Errors';

class UpdateName extends Component {
    constructor(props) {
        super(props);
        this.state = {
            studentID: '',
            courseCode: '',
            instructors: '',
            year: '',
            quarter: '',
            creditType: '',
            credits: '',
            workload: '',
            gradingTechniques: '',
            description: '',
            likedUsers: [],
            dislikedUsers: [],
            createdAt: '',
            editedAt:'',
            error: ''
        }
    }

    sendRequest = async (e) => {
        e.preventDefault();
        const { courseCode, instructors, year, quarter, creditType, credits, workload, gradingTechniques, description } = this.state;
        const sendData = { courseCode, instructors, year, quarter, creditType, credits, workload, gradingTechniques, description };
        const response = await fetch(api.base + api.handlers.evaluations, {
            method: "POST",
            body: JSON.stringify(sendData),
            headers: new Headers({
                "Authorization": localStorage.getItem("Authorization"),
                "Content-Type": "application/json"
            })
        });
        if (response.status >= 300) {
            const error = await response.text();
            console.log(error);
            this.setError(error);
            return;
        }
        alert("Evaluation added") // TODO make this better by refactoring errors
        const user = await response.json();
        this.props.setUser(user);
    }

    setValue = (e) => {
        this.setState({ [e.target.name]: e.target.value });
    }

    setError = (error) => {
        this.setState({ error })
    }

    render() {
        const { courseCode, instructors, year, quarter, creditType, credits, workload, gradingTechniques, description, error } = this.state;
        return <>
            <Errors error={error} setError={this.setError} />
            <div>Enter a new Evaluation</div>
            <form onSubmit={this.sendRequest}>
                <div>
                    <span>courseCode: </span>
                    <input name={"courseCode"} value={courseCode} onChange={this.setValue} />
                </div>
                <div>
                    <span>instructors: </span>
                    <input name={"instructors"} value={instructors} onChange={this.setValue} />
                </div>
                <div>
                    <span>year: </span>
                    <input name={"year"} value={year} onChange={this.setValue} />
                </div>
                <div>
                    <span>quarter: </span>
                    <input name={"quarter"} value={quarter} onChange={this.setValue} />
                </div>
                <div>
                    <span>creditType: </span>
                    <input name={"creditType"} value={creditType} onChange={this.setValue} />
                </div>
                <div>
                    <span>credits: </span>
                    <input name={"credits"} value={credits} onChange={this.setValue} />
                </div>
                <div>
                    <span>workload: </span>
                    <input name={"workload"} value={workload} onChange={this.setValue} />
                </div>
                <div>
                    <span>grading fairness: </span>
                    <input name={"gradingTechniques"} value={gradingTechniques} onChange={this.setValue} />
                </div>
                <div>
                    <span>description: </span>
                    <input name={"description"} value={description} onChange={this.setValue} />
                </div>

                <input type="submit" value="Add Evaluation" />
            </form>
        </>
    }

}

export default UpdateName;