import React from "react";
import AccountInfo from "./AccountInfo.tsx";
import { PostingInfoPage } from "./PostingInfo.tsx";
import { BusinessInfoPage } from "./BusinessInfo.tsx";
import { ApplicationInfo } from "./ApplicationInfo.tsx";
import { UserApplicationInfo } from "./UserApplicationInfo.tsx";

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

interface InfoPageProps {
    page: PageTag
    role: Role
}

export const InfoPage = ({ page, role }: InfoPageProps) => {
    switch (page) {
        case PageTag.Account:
            return (
                <AccountInfo />
            );
        case PageTag.Application:
            return (
                // (role === Role.Founder || role === Role.Admin) ? <ApplicationInfo /> : <UserApplicationInfo />
                <UserApplicationInfo />
            );
        case PageTag.Business:
            return (
                <BusinessInfoPage isAdmin={role === Role.Admin} />
            );
        case PageTag.Approvals:
            return (
                <AccountInfo />
            );
        case PageTag.Postings:
            return (
                <PostingInfoPage />
            );
        default:
            return (
                <AccountInfo />
            );
    }
}