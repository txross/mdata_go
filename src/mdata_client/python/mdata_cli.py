# Copyright 2017 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------------------

from __future__ import print_function

import argparse
import getpass
import logging
import os
import traceback
import sys
import pkg_resources
import re

from colorlog import ColoredFormatter

from client.mdata_client import MdClient
from client.mdata_exceptions import MdException


DISTRIBUTION_NAME = 'mdata-go'


DEFAULT_URL = 'http://127.0.0.1:8008'

ATTRIBUTE_PATTERN = r'^[^=]+={1}[^=]+$'

def create_console_handler(verbose_level):
    clog = logging.StreamHandler()
    formatter = ColoredFormatter(
        "%(log_color)s[%(asctime)s %(levelname)-8s%(module)s]%(reset)s "
        "%(white)s%(message)s",
        datefmt="%H:%M:%S",
        reset=True,
        log_colors={
            'DEBUG': 'cyan',
            'INFO': 'green',
            'WARNING': 'yellow',
            'ERROR': 'red',
            'CRITICAL': 'red',
        })

    clog.setFormatter(formatter)

    if verbose_level == 0:
        clog.setLevel(logging.WARN)
    elif verbose_level == 1:
        clog.setLevel(logging.INFO)
    else:
        clog.setLevel(logging.DEBUG)

    return clog


def setup_loggers(verbose_level):
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG)
    logger.addHandler(create_console_handler(verbose_level))

def key_value_pair(string):
    if not re.match(ATTRIBUTE_PATTERN, string):
        msg = "%s is not a key=value pair" % string
        raise argparse.ArgumentTypeError(msg)
    return string

def add_create_parser(subparsers, parent_parser):
    parser = subparsers.add_parser(
        'create',
        help='Creates a new product',
        description='Sends a transaction to create a new produc with the '
        'identifier <gtin>. This transaction will fail if the specified '
        'product already exists.',
        parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='unique GTIN-14 identifier for the new product')

    parser.add_argument(
        '--attributes',
        metavar="ATTR",
        nargs="*",
        default="",
        action="append",
        type=key_value_pair,
        help='provide key=value pair(s) of attributes')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--disable-client-validation',
        action='store_true',
        default=False,
        help='disable client validation')

    parser.add_argument(
        '--wait',
        nargs='?',
        const=sys.maxsize,
        type=int,
        help='set time, in seconds, to wait for product to commit')


def add_list_parser(subparsers, parent_parser):
    parser = subparsers.add_parser(
        'list',
        help='Displays information for all products',
        description='Displays information for all products in state, showing '
        'the product state for each product.',
        parents=[parent_parser])

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')


def add_show_parser(subparsers, parent_parser):
    parser = subparsers.add_parser(
        'show',
        help='Displays information about a product',
        description='Displays the product <gtin>, showing '
        'the product state',
        parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='identifier for the product')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')


def add_update_parser(subparsers, parent_parser):
    parser = subparsers.add_parser(
        'update',
        help='Update a product with new attributes',
        description='Sends a transaction to update attributes of a product '
        'with the identifier <gtin>. This transaction will fail if the '
        'specified product does not exist or if the attributes are malformed.',
        parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='identifier for the product')

    parser.add_argument(
        'attributes',
        metavar="ATTR",
        nargs="*",
        default="",
        action="append",
        type=key_value_pair,
        help='provide key=value pair(s) of attributes')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--wait',
        nargs='?',
        const=sys.maxsize,
        type=int,
        help='set time, in seconds, to wait for take transaction '
        'to commit')


def add_delete_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('delete', parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='Gtin identifier of the product to be deleted')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--wait',
        nargs='?',
        const=sys.maxsize,
        type=int,
        help='set time, in seconds, to wait for delete transaction to commit')

def add_deactivate_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('deactivate', parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='Gtin identifier of the product to be deactivated')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--wait',
        nargs='?',
        const=sys.maxsize,
        type=int,
        help='set time, in seconds, to wait for deactivate transaction to commit')

def add_activate_parser(subparsers, parent_parser):
    parser = subparsers.add_parser('activate', parents=[parent_parser])

    parser.add_argument(
        'gtin',
        type=str,
        help='Gtin identifier of the product to be activated')

    parser.add_argument(
        '--url',
        type=str,
        help='specify URL of REST API')

    parser.add_argument(
        '--username',
        type=str,
        help="identify name of user's private key file")

    parser.add_argument(
        '--key-dir',
        type=str,
        help="identify directory of user's private key file")

    parser.add_argument(
        '--auth-user',
        type=str,
        help='specify username for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--auth-password',
        type=str,
        help='specify password for authentication if REST API '
        'is using Basic Auth')

    parser.add_argument(
        '--wait',
        nargs='?',
        const=sys.maxsize,
        type=int,
        help='set time, in seconds, to wait for activate transaction to commit')


def create_parent_parser(prog_name):
    parent_parser = argparse.ArgumentParser(prog=prog_name, add_help=False)
    parent_parser.add_argument(
        '-v', '--verbose',
        action='count',
        help='enable more verbose output')

    try:
        version = pkg_resources.get_distribution(DISTRIBUTION_NAME).version
    except pkg_resources.DistributionNotFound:
        version = 'UNKNOWN'

    parent_parser.add_argument(
        '-V', '--version',
        action='version',
        version=(DISTRIBUTION_NAME + ' (Hyperledger Sawtooth) version {}')
        .format(version),
        help='display version information')

    return parent_parser


def create_parser(prog_name):
    parent_parser = create_parent_parser(prog_name)

    parser = argparse.ArgumentParser(
        description='Provides subcommands to create product master data '
        'by sending mdata transactions.',
        parents=[parent_parser])

    subparsers = parser.add_subparsers(title='subcommands', dest='command')

    subparsers.required = True

    add_create_parser(subparsers, parent_parser)
    add_list_parser(subparsers, parent_parser)
    add_show_parser(subparsers, parent_parser)
    add_update_parser(subparsers, parent_parser)
    add_delete_parser(subparsers, parent_parser)
    add_deactivate_parser(subparsers, parent_parser)
    add_activate_parser(subparsers, parent_parser)

    return parser

def data_splitter(product_data):
    gtin = product_data[0]
    attributes = product_data[1:len(product_data)-1]
    state = product_data[len(product_data)-1]
    return gtin, attributes, state

def do_list(args):
    url = _get_url(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=None)

    product_list = [
        product.split(',')
        for products in client.list(auth_user=auth_user,
                                 auth_password=auth_password)
        for product in products.decode().split('|')
    ]

    if product_list is not None:
        fmt = "%-15s %-15.15s %-15.15s %-9s %s"
        print(fmt % ('GTIN', 'ATTRIBUTES', 'STATE'))
        for product_data in product_list:

            gtin, attributes, state = data_splitter(product_data)

            print(fmt % (gtin, attributes, state))
    else:
        raise MdException("Could not retrieve product listing.")


def do_show(args):
    gtin = args.gtin

    url = _get_url(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=None)

    data = client.show(gtin, auth_user=auth_user, auth_password=auth_password)

    if data is not None:

        gtin, attributes, state = {
            gtin: (gtin, attributes, state)
            for gtin, attributes, state in [
                data_splitter(product.split(','))
                for product in data.decode().split('|')
            ]
        }[gtin]

        print("product:     : {}".format(gtin))
        print("ATTRIBUTES  : {}".format(attributes))
        print("STATE     : {}".format(state))
        print("")

    else:
        raise MdException("product not found: {}".format(gtin))


def do_create(args):
    gtin = args.gtin

    url = _get_url(args)
    keyfile = _get_keyfile(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=keyfile)

    if args.wait and args.wait > 0:
        response = client.create(
            gtin, wait=args.wait,
            auth_user=auth_user,
            auth_password=auth_password)
    else:
        response = client.create(
            gtin, auth_user=auth_user,
            auth_password=auth_password)

    print("Response: {}".format(response))


def do_update(args):
    gtin = args.gtin
    attributes = args.attributes

    url = _get_url(args)
    keyfile = _get_keyfile(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=keyfile)

    if args.wait and args.wait > 0:
        response = client.update(
            gtin, attributes, wait=args.wait,
            auth_user=auth_user,
            auth_password=auth_password)
    else:
        response = client.update(
            gtin, attributes,
            auth_user=auth_user,
            auth_password=auth_password)

    print("Response: {}".format(response))


def do_delete(args):
    gtin = args.gtin

    url = _get_url(args)
    keyfile = _get_keyfile(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=keyfile)

    if args.wait and args.wait > 0:
        response = client.delete(
            gtin, wait=args.wait,
            auth_user=auth_user,
            auth_password=auth_password)
    else:
        response = client.delete(
            gtin, auth_user=auth_user,
            auth_password=auth_password)

    print("Response: {}".format(response))

def do_set_state(args, state):
    gtin = args.gtin

    url = _get_url(args)
    keyfile = _get_keyfile(args)
    auth_user, auth_password = _get_auth_info(args)

    client = MdClient(base_url=url, keyfile=keyfile)

    if args.wait and args.wait > 0:
        response = client.set_state(
            gtin, state, wait=args.wait,
            auth_user=auth_user,
            auth_password=auth_password)
    else:
        response = client.set_state(
            gtin, state, auth_user=auth_user,
            auth_password=auth_password)

    print("Response: {}".format(response))

def _get_url(args):
    return DEFAULT_URL if args.url is None else args.url


def _get_keyfile(args):
    username = getpass.getuser() if args.username is None else args.username
    home = os.path.expanduser("~")
    key_dir = os.path.join(home, ".sawtooth", "keys")

    return '{}/{}.priv'.format(key_dir, username)


def _get_auth_info(args):
    auth_user = args.auth_user
    auth_password = args.auth_password
    if auth_user is not None and auth_password is None:
        auth_password = getpass.getpass(prompt="Auth Password: ")

    return auth_user, auth_password


def main(prog_name=os.path.basename(sys.argv[0]), args=None):
    if args is None:
        args = sys.argv[1:]
    parser = create_parser(prog_name)
    args = parser.parse_args(args)

    if args.verbose is None:
        verbose_level = 0
    else:
        verbose_level = args.verbose

    setup_loggers(verbose_level=verbose_level)
    
    if args.command == 'create':
        do_create(args)
    elif args.command == 'list':
        do_list(args)
    elif args.command == 'show':
        do_show(args)
    elif args.command == 'update':
        do_update(args)
    elif args.command == 'delete':
        do_delete(args)
    elif args.command == 'activate':
        do_set_state(args, "ACTIVE")
    elif args.command == 'deactivate':
        do_set_state(args, "INACTIVE")
    else:
        raise MdException("invalid command: {}".format(args.command))


def main_wrapper():
    try:
        main()
    except MdException as err:
        print("Error: {}".format(err), file=sys.stderr)
        sys.exit(1)
    except KeyboardInterrupt:
        pass
    except SystemExit as err:
        raise err
    except BaseException as err:
        traceback.print_exc(file=sys.stderr)
        sys.exit(1)
