export default {
    base: "https://apicourseeval.info441-deploy.me",
    testbase: "https://localhost:4000",
    handlers: {
        users: "/v1/users",
        myuser: "/v1/users/me",
        myuserAvatar: "/v1/users/me/avatar",
        sessions: "/v1/sessions",
        sessionsMine: "/v1/sessions/mine",
        resetPasscode: "/v1/resetcodes",
        passwords: "/v1/passwords/",
        courses: "/v1/courses",
        evaluations: "/v1/evaluations/"
    }
}