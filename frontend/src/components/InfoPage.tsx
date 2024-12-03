import React from "react";
import AccountInfo from "./AccountInfo.tsx";
import { PostingInfoPage } from "./PostingInfo.tsx";
import { BusinessInfoPage } from "./BusinessInfo.tsx";
import { ApplicationInfo } from "./ApplicationInfo.tsx";

export enum PageTag {
    Account = "account",
    Application = "application",
    Business = "business",
    Approvals = "approvals",
    Postings = "postings"
}

export enum Role {
    User = "user",
    Admin = "admin"
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
                <ApplicationInfo />
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