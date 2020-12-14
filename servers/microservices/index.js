const mongoose = require('mongoose');
const express = require('express');
const { course, evaluation } = require('./schemas');

const mongoEndpoint = "mongodb://customMongoContainer:27017/rabbit";
const port = 80;

const Course = mongoose.model("Course", course);
const Evaluation = mongoose.model("Evaluation", evaluation);

// set up express
const app = express();
app.use(express.json());

// A function to connect to the mongo endpoint, used for refreshing on disconnect.
const connect = () => {
	mongoose.connect(mongoEndpoint);
}

// create & add default courses
defaultCourses();

// app.all('*',function(req,res,next)
// {
//     if (!req.get('Origin')) return next();

//     res.set('Access-Control-Allow-Origin','*');
//     res.set('Access-Control-Allow-Methods','GET,POST,DELETE,PATCH');
//     res.set('Access-Control-Allow-Headers','X-Requested-With,Content-Type');

//     if ('OPTIONS' == req.method) return res.send(200);

//     next();
// });

// get all courses or courses that match search param
app.get("/v1/courses", async (req, res) => {
	// 200: Successful response with all course information
	// 500: Internal server error 

	try {
		// get search query 
		const query = req.query['code'] ? req.query['code'] : null;

		if (query == null) {
			// get all courses
			const courses = await Course.find();
			res.set('Content-Type', 'application/json');
			res.status(200).json(courses);
		} else {
			// get courses with matching code 
			const matchingCourses = await Course.find({code: {"$regex": query, "$options": "i"}})
			
			// return courses 
			res.set ("Content-Type", "application/json");
			res.status(200).json(matchingCourses);
		}
		
	} catch {
		res.status(500).send("There was an issue getting courses");
	}
});

// get specific course based on given course ID 
app.get("/v1/courses/:courseID", async (req, res) => {
    // 200: Successful response with course information 
	// 500: Internal server error
	
    try {
        // get course with id from req
        const courseID = req.params['courseID'];
        const course = await Course.find({ id: courseID });

        res.set("Content-Type", "application/json");
        res.status(200).json(course);
    } catch {
        res.status(500).send("There was an issue getting the course");
    }
});

// get all evaluations for that course
app.get("/v1/courses/:courseID/evaluations", async (req, res) => {
	// 200: Successful response with all evals in that course
	// 500: Internal server error 
	try {
		const courseID = req.params['courseID'];

		const evaluations = await Evaluation.find({ courseID: courseID });

        res.set("Content-Type", "application/json");
        res.status(200).json(evaluations);
	} catch {
		res.status(500).send("There was an issue getting evaluations for this course");
	}
});

//add new evaluation 
app.post("/v1/evaluations/", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}

    const { courseCode, instructors, year, quarter, creditType, credits, workload, gradingTechniques, description } = req.body;
	const course = await Course.find({ code: courseCode });

	//to get studentID
	var usr = JSON.parse(XUser);
	//get number of documents in evaluation
	const Lastid = await Evaluation.countDocuments({});
	const id = Lastid+1;
    const createdAt = new Date();
	const evaluation = { 
		id: id,
		studentID: usr.id,
        courseID: course.id,
        instructors: instructors,
        year: year,
        quarter: quarter,
        creditType: creditType,
        credits: credits,
        workload: workload,
        gradingTechniques: gradingTechniques,
		description: description,
        createdAt: createdAt
	};

	const query = new Evaluation(evaluation);
	query.save((err, newEvaluation) => {
		if (err) {
			console.log(err);
			res.status(500).send("Unable to create new evaluation");
			return;
        }
        
        res.set("Content-Type", "application/json");
        res.status(201).json(newEvaluation);
    });
});

//update evaluation text description based on specific evaluation ID
app.patch("/v1/evaluations/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get evaluation id 
	const evaluationId = req.params.id;
	const specificEvaluation = await Evaluation.find({ id: evaluationId });
	//get user
	var usr = JSON.parse(XUser);
	//if user not creator of Evaluation
	if (specificEvaluation[0].studentID != usr.id) {
		res.status(403).send("Forbidden User");
		return;
	}
	// get new description 
	const {description} = req.body;
	// update members
	await Evaluation.where({ id: evaluationId }).updateOne({ description: description });
	const updatedEvaluation = await Evaluation.find({ id: evaluationId });
	res.set("Content-Type", "application/json");
	res.json(updatedEvaluation);

});

app.delete("/v1/evaluations/:id", async (req, res) => {
	// verify user authorization
	var XUser = req.header('X-User');
	if(!XUser){
		res.status(401).send("User Unauthorized");
		return;
	}
	//get evaluation id 
	const evaluationId = req.params.id;
	const specificEvaluation = await Evaluation.find({ id: evaluationId });
	//get user
	var usr = JSON.parse(XUser);
	//if user not creator of Evaluation
	if (specificEvaluation[0].studentID != usr.id) {
		res.status(403).send("Forbidden User");
		return;
	}
	// delete Evaluation
    await Evaluation.deleteOne({ id: evaluationId });
    res.set("Content-Type", "text/plain");
	res.send("Successfullly deleted evaluation");
});

connect();
mongoose.connection.on('error', console.error)
	.on('disconnected', connect)
	.once('open', main);

async function main() {
	app.listen(port, "", () => {
		console.log(`Server listening on port ${port}`);
	});
}

function defaultCourses() {
	// create default courses
	const info441 = { 
		id: 1,
		code: 'INFO441',
		title: 'Server-side Development',
		description: 'Introduces server-side web development programming, services, tools, protocols, \
		best practices and techniques for implementing data-driven and scalable web applications. \
		Connects topics from human-centered design, information architecture, databases, data analytics and security to build a solution.',
		credits: 5,
	};
	query = new Course(info441);
	query.save();

	const psych101 = { 
		id: 2,
		code: 'PSYCH101',
		title: 'Introduction to Psychology',
		description: 'Surveys major areas of psychological science. Core topics include human social behavior, personality, psychological disorders and treatment, learning, memory, human development, biological influences, and research methods. Related topics may include sensation, perception, states of consciousness, thinking, intelligence, language, motivation, emotion, stress and health, cross-cultural psychology, and applied psychology.',
		credits: 5,
	};
	query = new Course(psych101);
	query.save();

	const cse143 = { 
		id: 3,
		code: 'CSE143',
		title: 'Computer Programing II',
		description: 'Continuation of CSE 142. Concepts of data abstraction and encapsulation including stacks, queues, linked lists, binary trees, recursion, instruction to complexity and use of predefined collection classes',
		credits: 5,
	};
	query = new Course(cse143);
	query.save();

	const info200 = { 
		id: 4,
		code: 'INFO200',
		title: 'Computer Programing II',
		description: 'Information as an object of study, including theories, concepts, and principles of information, information seeking, cognitive processing, knowledge representation and restructuring, and their relationships to physical and intellectual access to information. Development of information systems for storage, organization, and retrieval. Experience in the application of theories, concepts, and principles.',
		credits: 5,
	};
	query = new Course(info200);
	query.save();

	const math308 = {
		id: 5,
		code: 'MATH308',
		title: 'Matrix Algebra With Applications',
		description: 'Systems of linear equations, vector spaces, matrices, subspaces, orthogonality, least squares, eigenvalues, eigenvectors, applications. For students in engineering, mathematics, and the sciences.',
		credits: 3,
	}
	query = new Course(math308);
	query.save();

	const ess100 = {
		id: 6,
		code: 'ESS100',
		title: 'Dinosaurs',
		description: 'Biology, behavior, ecology, evolution, and extinction of dinosaurs, and a history of their exploration. With dinosaurs as focal point, course also introduces the student to how hypotheses in geological and paleobiological science are formulated and tested.',
		credits: 2,
	}
	query = new Course(ess100);
	query.save();

	const educ251 = {
		id: 7,
		code: 'EDUC251',
		title: 'Seeking Educational Equity And Diversity',
		description: 'Introduces the need for and challenges in establishing educational equity and diversity. Discussions explore theories, historical trends, and ongoing debates. Readings draw from academic and popular sources, and class sessions include use of multimedia resources and experiential activities.',
		credits: 5,
	}
	query = new Course(educ251);
	query.save();

	const ling200 = {
		id: 8,
		code: 'LING200',
		title: 'Introduction To Linguistics',
		description: 'Language as the fundamental characteristic of the human species; diversity and complexity of human languages; phonological and grammatical analysis; dimensions of language use; and language acquisition and historical language change.',
		credits: 5,
	}
	query = new Course(ling200);
	query.save();

	const testEval = {
		id: 1,
		studentID: 1,
		courseID: 1,
		instructors: [{"name":"Kyle Thayer"}],
		year: new Date().getYear(),
		quarter: 'AU20',
		creditType: 'I&S',
		credits: 5,
		workload: 10,
		gradingTechniques: 3,
		description: 'Most difficult Info class',
		likedUsers: [{"studentID":1}],
		dislikedUsers: [{"studentID":2}],
		createdAt: new Date(),
		editedAt: new Date()
	}
	query = new Evaluation(testEval);
	query.save();
}
