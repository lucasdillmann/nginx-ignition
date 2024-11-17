import {StyleBook} from "../../contracts/StyleBook";

export default {
    container: {
        position: "absolute",
        width: "100%",
        height: "100%",
        top: 0,
        left: 0,
    },
    header: {
        padding: 0,
        background: "#FFF",
    },
    title: {
        padding: "30px 20px 0",
        margin: 0,
    },
    content: {
        margin: "24px 16px",
        padding: 24,
        minHeight: 280,
        background: "#FFF",
        borderRadius: 4,
    },
    toggleButton: {
        fontSize: "16px",
        width: 64,
        height: 64,
    },
    logo: {
        color: "#FFF",
        fontSize: 18,
        padding: "30px 20px 20px",
    },
    logoLink: {
        color: "#FFF",
        fontSize: 22,
    },
    menuItem: {
        margin: "5px 10px",
        padding: 10,
        width: 230,
        borderRadius: 4,
    },
} as StyleBook
