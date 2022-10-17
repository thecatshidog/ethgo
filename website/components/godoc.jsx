
import Link from 'next/link'

const goDocRef = "https://pkg.go.dev/github.com/thecatshidog/ethgo/"

export default function GodocLink({children, href}) {
    return <Link href={`${goDocRef}${href}`}>{children}</Link>
}
