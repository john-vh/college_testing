import React from "react";
import AccountInfo from "./AccountInfo.tsx";
import { PostingInfoPage } from "./PostingInfo.tsx";
import { BusinessInfoPage } from "./BusinessInfo.tsx";
import { ApplicationInfo } from "./ApplicationInfo.tsx";
import { UserApplicationInfoPage } from "./UserApplicationInfo.tsx";
import { Route, Routes, useParams } from "react-router-dom";
import { useGetRole } from "../hooks/useAccountInfo.ts";

export enum PageTag {
    Account = "account",
    Application = "application",
    Business = "business",
    Approvals = "approvals",
    Postings = "postings"
}

export enum Role {
    User = "user",
    Admin = "admin",
    Founder = "founder"
}

export const InfoPage = () => {
    const { page } = useParams<{ page: string }>();
    const role = useGetRole();

    switch (page) {
        case "account":
            return <AccountInfo />;
        case "application":
            return (
                (role === Role.Founder || role === Role.Admin) ? <ApplicationInfo isAdmin={role === Role.Admin} /> : <UserApplicationInfoPage />
            );
        case "business":
            return <BusinessInfoPage isAdmin={role === Role.Admin} />;
        case "approvals":
            return <AccountInfo />;
        case "postings":
            return <PostingInfoPage isAdmin={role === Role.Admin} />;
        default:
            return <AccountInfo />;
    }
}